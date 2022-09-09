# Sodor
Sodor is a distributed and extensible scheduler system, with lower operating expenses and high performance.

# Some Expected Usage
* fat_ctrl
  * ./fat_ctrl --port xxxx
* thomas
  * ./thomas --grpc.port xxxx 
  * run job
    * ./thomas --standalone run_job --job.workdir ./$jobid/
    * job.workdir
      * meta.json 
        * write by thomas
        * include
          * process_id
          * ttl
      * status.json
        * write by thomas-job
        * include
          * total time
          * exit code
          * exit msg
          * process id
      * log
        * task1.log
        * task2.log
      * data
