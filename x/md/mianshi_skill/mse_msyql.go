package mianshiskill

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type BankService struct {
	db *sql.DB
}

// 使用 READ COMMITTED 隔离级别

// 通过 FOR UPDATE 锁定相关账户行

// 采用乐观锁机制（version字段）
// 按帐户ID顺序更新（预防死锁）排序消除环路
// 场景1：交叉更新
// 问题描述
// 事务A更新账户1→账户2，事务B更新账户2→账户1

// 解决方案
// 统一按账户ID顺序更新

func (s *BankService) Transfer(ctx context.Context, fromAccount, toAccount string, amount int) error {

	attempt := 0
	maxRetries := 3
	lockTimeout := 10 * time.Second
	var lastErr error
	for attempt < maxRetries {
		ctx, cancel := context.WithTimeout(ctx, lockTimeout)
		defer cancel()
		tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			return err
		}

		err = s.transfer(ctx, tx, fromAccount, toAccount, amount)
		if err == nil {
			if commitErr := tx.Commit(); commitErr == nil {
				return nil
			} else {
				lastErr = commitErr
				//attempt++
			}
		} else {
			lastErr = err
		}

		// 最后回滚
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			lastErr = rollBackErr
			break
			// 记录日志回滚失败
		}
		//判断是否需要重试
		if ifNeedRetry(lastErr) {
			//attempt++
			continue
		}
		break
	}

	// defer tx.Rollback()

	// 或者先排序帐户
	return lastErr
}

func (s *BankService) transfer(ctx context.Context, tx *sql.Tx, fromAccount, toAccount string, amount int) error {

	var fromBalance float64
	var fromVersion int

	err := tx.QueryRowContext(ctx, "SELECT balance, version FROM accounts WHERE account_id =? FOR UPDATE", fromAccount).Scan(&fromBalance, &fromVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("from account not found")
		}
		return err
	}
	if fromBalance < float64(amount) {
		return fmt.Errorf("insufficient balance")
	}

	var toVersion int
	err = tx.QueryRowContext(ctx, "SELECT version FROM accounts WHERE account_id =? FOR UPDATE", toAccount).Scan(&toVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("to account not found")
		}
		return err
	}

	accounts := sortAccounts(fromAccount, toAccount)

	for _, acc := range accounts {
		// 按帐户ID顺序更新（预防死锁）,先更新小的
		switch acc {
		case fromAccount:
			_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ?, version=version+1 WHERE account_id = ? and version=?", amount, fromAccount, fromVersion)
		case toAccount:
			_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance +? version=version+1 WHERE account_id =? and version=?", amount, toAccount, toVersion)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func ifNeedRetry(err error) bool {
	return err != nil && err.Error() == "deadlock detected"
}

func sortAccounts(ids ...string) []string {
	if ids[0] > ids[1] {
		return []string{ids[1], ids[0]}
	}
	return ids
}
