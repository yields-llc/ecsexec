# ecsexec
"ecsexec" is a command that makes "aws ecs execute-command" useful.  
"ecsexec" does not require a task ID. Instead, service name is required.  

"aws ecs execute-command" usage example is as follows.

```bash
$ aws ecs execute-command \
    --cluster MyCluster \
    --task arn:aws:ecs:us-east-1:123456789012:task/MyCluster/d789e94343414c25b9f6bd59eEXAMPLE \
    --container MyContainer \
    --interactive \
    --command "/bin/sh"
```

"--task" is the task ID, but this task ID is tricky.  
The task ID is temporary and not something the user can guess. 

"ecsexec" usage example is as follows.

```bash
$ ecsexec exec \
    -profile <AWS_PROFILE> \
    -cluster MyCluster \
    -service MyService \
    -container MyContainer \
    -command "/bin/sh"
```

"-service" is the service name, but service name is not tricky.  
Service names are constant and user-friendly.

"ecsexec" simplifies command execution by specifying the service name rather than the task ID.

"ecsexec" will start the service with desired count 1 if the service is not running, and after execution, "ecsexec" will stopp the service with desired count 0.  

```bash
$ ecsexec exec -profile=yields -cluster=demo -service=demo -container=demo -command="ls -la /var/log"

"demo" service was stopped. it will be started now...
"demo" service has been started.
waiting for SSM agent running...
SSM agent is running


The Session Manager plugin was installed successfully. Use the AWS CLI to start a session.


Starting session with SessionId: ecs-execute-command-08e188379c544c240
total 16
drwxr-xr-x    1 root     root          4096 Sep 17 16:22 .
drwxr-xr-x    1 root     root          4096 Aug  7 13:09 ..
drwx------    3 root     root          4096 Sep 17 16:22 amazon


Exiting session with sessionId: ecs-execute-command-08e188379c544c240.


"demo" service was started. it will be stopped now...
"demo" service has been stopped.
```

If the service is started, execute the command by using the one task that is running.

```bash
$ ecsexec exec -profile=yields -cluster=demo -service=demo -container=demo -command="ls -la /var/log"

"demo" service was already started (running count = 1).


The Session Manager plugin was installed successfully. Use the AWS CLI to start a session.


Starting session with SessionId: ecs-execute-command-0d685051f95db36a6
total 16
drwxr-xr-x    1 root     root          4096 Sep 17 16:27 .
drwxr-xr-x    1 root     root          4096 Aug  7 13:09 ..
drwx------    3 root     root          4096 Sep 17 16:27 amazon


Exiting session with sessionId: ecs-execute-command-0d685051f95db36a6.
```

## Required
"ecsexec" requires the following programs to be installed.  

- Install [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- Install [Session Manager plugin for the AWS CLI](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

## Restrictions
If you find any problems, please create an issue or pull-request.

- Not tested on Windows

## Install
```bash
$ go install github.com/yields-llc/ecsexec/cmd/ecsexec
```
