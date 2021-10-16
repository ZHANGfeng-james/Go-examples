Go 语言 sync 标准库提供了**基本的并发原语**，比如：互斥锁 sync.Mutex。另外的一些提供并发操作的还有：

* Once 和 WaitGroup
* Channel 机制

特别提醒：sync 库中定义的**值**，**不允许拷贝**。



### 1 结构体

~~~bash
ant@MacBook-Pro ~ % go doc sync |grep "^type"
type Cond struct{ ... }
type Map struct{ ... }
type Mutex struct{ ... }
type Once struct{ ... }
type Pool struct{ ... }
type RWMutex struct{ ... }
type WaitGroup struct{ ... }

type Locker interface{ ... }
~~~

依据上述，大致梳理出如下实体：

**One：sync.Mutex**——













**Two：sync.RWMutex**——













**Three：sync.Cond**——













**Four：sync.Map**——











**Five：sync.Once**——











**Six：sync.Pool**——









**Seven：sync.WaitGroup**——



