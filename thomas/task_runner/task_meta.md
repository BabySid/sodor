# task runner meta file
* thomas/
  * tasks/
    * $job_id/
      * $task_id/
        * $job_instance_$task_instance/
          * task_request.json
          * task_response.json
          * data/
          * logs/
        * $job_instance_$task_instance/
          * ...
  * status
    * thomas.sqlite
* another design
  * thomas accept run-task-request
    * fork thomas to process the request
    * key point is 
      * the child process the request and response to fat_ctrl
        * pros
          * 对于thomas异常维护时容错性较好，子进程可以独立处理请求
          * 子进程处理完毕可以高时效性和fat-ctrl交互
        * cons
          * 子进程一旦异常（如被kill），fat-ctrl无法感知。需要thomas重新执行
          * thomas仍需要管理子进程的逻辑
          * fat-ctrl只能依赖任务超时来check 子进程的状态
      * the parent process the request and response to fat_ctrl
        * pros
          * 总体稳定性好，thomas作为强绑定、一级公民，其生命周期较子进程强，因此可以管理子进程的生命周期（如杀死、重启等）
          * 避免众多子进程对fat-ctrl的压力
          * 综合收集子进程的执行情况，形成thomas的负载评估参数
        * cons
          * 当thomas异常时，对于子进程的感知依赖于thomas的恢复进度
          * thomas重启后需要load原来子进程的状态，并轮询检测（时效性略慢几秒）
      * 最终：选择后者。
        * 核心是因为对于实际的众多子进程，在fat-ctrl和其之间，需要一个管理的服务。损失的时效性可忽略不计。
        * 异常时，本身对于thomas所在的节点已经需要修复了。其上的任务甚至进一步可以按照节点宕机处理。所以子进程自身的容错略高的优势不明显