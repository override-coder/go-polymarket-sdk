# Relayer Context 入口

旧方法保持签名及默认行为，委托使用 context.Background() 的新方法。
需要业务取消/超时时，使用以下同名加 Context 后缀的方法，并在第一个参数传入 ctx：

- GetNonceContext、GetDeployedContext、GetTransactionsContext
- DeployContext、BuildDeployTxContext、BuildTxContext
- ExecuteContext、ExecuteByTxContext
- DeployDepositWalletContext、BuildDepositWalletBatchContext
- ExecuteDepositWalletBatchContext、ExecuteDepositWalletBatchByTxContext

既有 GetTransactionContext、PollUntilState 继续可用。
Context 贯穿 deployed/nonce 查询、签名前检查、HTTP 提交及重试等待。
已取消的调用在签名前退出；若在自定义签名回调执行时取消，回调结束后不再发出提交请求。
现有 SignatureFunc 不接收 context，SDK 无法强制中断调用方的阻塞签名回调。

取消不能撤销已经被服务端接受或上链的交易。提交阶段返回取消/网络错误时，
仍须用交易记录、Relayer ID 和回执核实结果，不能直接假定未提交或盲目重试。

GetTransactions / GetTransactionsContext 配置了 Builder 凭据时，签名失败直接返回 error，
不再静默退化为无认证请求。完全未配置 Builder 凭据时保留原先查询行为；
这不保证服务端允许匿名访问。已签名 ByTx 入口仍不要求本地签名函数。
