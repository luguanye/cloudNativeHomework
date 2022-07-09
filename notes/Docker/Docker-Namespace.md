## Linux Namespace

[toc]

Linux Namespace 是一种 Linux 内核提供的资源隔离方案：

- 系统可以为进程分配不同的 Namespace；
- 并保证不同的 Namespace 资源独立分配、进程彼此隔离；

#### 进程数据结构

```c
struct task_struct {
  ...
  /* namespaces */
  struct nsproxy *nsproxy;
  ...
}
```

#### namespace 数据结构

```c
struct nsproxy {
  atomic_t count;
  struct uts_namespace *uts_ns;
  struct ipc_namespace *ipc_ns;
  struct mnt_namespace *mnt_ns;
  struct pid_namespace *pid_ns_for_children;
  struct net *net_ns;
}
```

#### Linux 内核的6种 namespace

| Namespace | 系统调用参数  | 隔离内容                   |
| :-------- | :------------ | :------------------------- |
| UTS       | CLONE_NEWUTS  | 主机名与域名               |
| IPC       | CLONE_NEWIPC  | 信息量、消息队列和共享内存 |
| PID       | CLONE_NEWPID  | 进程编号                   |
| Network   | CLONE_NEWNET  | 网络设备、网络栈、端口等等 |
| Mount     | CLONE_NEWNS   | 挂载点（文件系统）         |
| User      | CLONE_NEWUSER | 用户和用户组               |

- uts namespace

  UTS（UNIX Time-sharing System）namespace 允许每个container 拥有独立的 hostname 和 domain name，在网络上被视为一个独立的节点而非 host 上的一个进程。

- ipc namespace

  - 容器中进程交互还是采用 linux 常见的交互（ipc - inter process communication），包括信号量、共享内存、消息队列。
  - 容器的进程交互实际上还是 host 上具有相同 pid namespace 的进程间交互，因此需要再 ipc 资源申请时加入 namespace 信息 - 每个 ipc 资源有一个唯一的 32 位 id。

- pid namespace

  每个 namespace 中 pid相互隔离。

- net namespace

  - 每个 net namespace 有独立的 network device、ip address、ip routing tables、/proc/net 目录。
  - Docker 默认采用 veth 的方式将容器中的虚拟网卡同 host 上的一个 docker bridge：docker0 连接在一起。

- mnt namespace

  mnt namespace 允许不同 namespace 的进程看到的文件结构不同。

- user namespace

  每个 container 可以有不同的 user 和 group id，container 内部用container 内部的用户执行程序，而不是 host 的用户。

## namespace 系统调用

#### clone（创建一个新进程并放入到新的 namespace 中）

```c
int clone(int (*child_func)(void *), void *child_stack, int flags, void *arg);
```

通过 flags 指定新建的 namespace 类型：CLONE_NEWIPC、CLONE_NEWNS、CLONE_NEWNET、CLONE_NEWPID、CLONE_NEWUSER和CLONE_NEWUTS

#### setns（调用将当前进程加入到已有的 namespace 中）

```c
int setns(int fd, int nstype);
```

- 参数fd表示指向 /proc/[pid]/ns/ 目录里相应 namespace 对应的文件，表示要加入哪个 namespace 。
- 参数 nstype 让调用者可以去检查 fd 指向的 namespace 类型是否符合我们实际的要求，如果填0表示不检查。

#### unshare（使当前进程退出指定的 namespace 并加入到新创建的namespace中）

```c
int unshare(int flags);
```

- unshare() 运行在原先的进程上，不需要启动一个新进程
- flags 指定一个或者多个上面的CLONE_NEW*， 这样**当前进程**就退出当前指定类型的 namespace 并加入新创建的 namespace。

## namespace 常用操作

##### 查看当前系统的 namespace

```shell
lsns -t<type>
```

##### 查看某进程的 namespace

```
ls -la /proc/<pid>/ns/
```

##### 进入 namespace 运行命令

```
nsenter -t <pid> -n ip addr
```

##### 查看/proc/[pid]/ns文件

```shell
 ls -l /proc/2597/ns
total 0
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 cgroup -> cgroup:[4026531835]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 ipc -> ipc:[4026531839]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 mnt -> mnt:[4026531840]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 net -> net:[4026531957]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 pid -> pid:[4026531836]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 user -> user:[4026531837]
lrwxrwxrwx 1 zhangxa zhangxa 0 Mar 2 06:42 uts -> uts:[4026531838]
```

