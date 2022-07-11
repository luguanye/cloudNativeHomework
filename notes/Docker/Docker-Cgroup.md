Cgroups 在不同的系统资源管理子系统中以层级树（hierachy）的方式来组织管理，每个 cgroup 都可以包含子 cgroup，因此子 cgroup 的资源除了受本 cgroup 影响外还受父 cgroup 影响

## Linux Cgroup

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

#### 

