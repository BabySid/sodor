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