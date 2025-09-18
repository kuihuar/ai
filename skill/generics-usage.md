### Go 泛型：实际业务的推荐使用场景

- **通用容器与数据结构**
  - Set/Map 包装：`Set[T comparable]`、`MultiMap[K comparable, V any]`
  - 线性结构：`Stack[T]`、`Queue[T]`、`Ring[T]`
  - 树/堆/图：`Heap[T]`、`BST[T]`
  - 缓存：`LRUCache[K comparable, V any]`、`TTLCache[K,V]`

- **通用算法与工具函数**
  - 排序/去重/集合运算：`SortBy[T]`、`Unique[T comparable]`、`Intersect[T comparable]`
  - 数值：`Sum[T constraints.Integer|constraints.Float]`、`Min/Max[T constraints.Ordered]`
  - 函数式：`Map[T,R]`、`Filter[T]`、`Reduce[T,R]`
  - 迭代/流：`Iterator[T]`、`Stream[T]`

- **基础设施与中间件**
  - 重试/熔断/限流：`Retry[T any](op func() (T,error))`、`CircuitBreaker[T]`
  - 池化：`Pool[T any]`（连接/对象）
  - 事件/消息：`EventBus[T any]`、`PubSub[T any]`

- **编解码与校验**
  - `Decode[T any]([]byte) (T,error)`、`Encode[T any](v T) []byte`
  - `Validate[T any](v T) error`、统一壳 `Envelope/Result[T any]`

- **数据访问（Repository/DAO）**
  - `Repository[T any, ID comparable]`（CRUD、分页）
  - 查询构造器：`QueryBuilder[T]`、分页 `Page[T]`

- **HTTP/RPC 客户端**
  - `DoJSON[Req,Resp any](ctx, url, req) (Resp,error)`
  - 中间件链：`Middleware[T any]`

---

- **不建议使用的场景**
  - 强领域特定、仅服务单一实体的逻辑
  - 规则复杂、依赖接口多态的演进性需求
  - 抽象成本高、可读性下降或已知具体类型更高效

- **经验法则**
  - 同一套逻辑要服务多种类型，且能用约束清晰表达共性时用泛型
  - 没有稳定抽象共性、只是“看起来类似”，不要强行使用泛型 