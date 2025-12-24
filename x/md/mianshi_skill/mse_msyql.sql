create table accounts(
    accountNumber varchar(20) primary key,
    balance decimal(15, 2) not null,
    version int not null default 0
)engine=innodb default charset=utf8mb4;
insert into accounts(accountNumber, balance)values('1234567890', 1000.00), ('0987654321', 1000.00);
delimiter $$;
create procedure transferFunds(
    in fromAccount varchar(20),
    in toAccount varchar(20),
    in amount decimal(15, 2)
);

begin
declare fromBalance decimal(15, 2);
declare fromVersion int;
declare toVersion int;

start transaction;

-- 检查源账户是否存在且余额充足
select balance,version 
into fromBalance,fromVersion
from accounts
where accountNumber = fromAccount 
for update;

if row_count() = 0 then
    rollback;
    signal sqlstate '45000' set message_text = 'Source account not found';
end if;

-- 检查目标账户是否存在
select version
into toVersion
from accounts
where accountNumber = toAccount
for update;
if row_count() = 0 then
    rollback;
    signal sqlstate '45000' set message_text = 'Target account not found';
end if;

-- 检查源账户余额是否充足
if fromBalance < amount then
    rollback;
    signal sqlstate '45000'
    set message_text = 'Insufficient funds';
end if;


-- 更新转出账户余额带乐观锁
update accounts
set balance = balance - amount,version = version + 1
where accountNumber = fromAccount and version = fromVersion;
if row_count() = 0 then
    rollback;
    signal sqlstate '45000'
    set message_text = 'Optimistic locking failed  for source account';
end if;

-- 更新转入账户余额
update accounts
set balance = balance + amount, version = version + 1
where accountNumber = toAccount and version=toVersion;

if row_count() = 0 then
    rollback;
    signal sqlstate '45000'
    set message_text = 'Optimistic locking failed for target account';
end if;

commit;
end $$ 
delimiter;
transferFunds('1234567890', '0987654321', 100.00);