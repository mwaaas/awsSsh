**_awsSsh_**

This is a utility to help us ssh to aws elastic beanstalk and 
ecs 

There is already aws tool `eb ssh` which works well, awsSsh is not meant
to replace it. It's only that for `eb ssh` to work you need to have done 
`eb init` which works well but if you don't want to do that this tool is for you.


*`Running`*
`Export your aws credentials`

    export AWS_ACCESS_KEY_ID={your_key}; export AWS_SECRET_ACCESS_KEY={your_secret_key}; export AWS_REGION=eu-west-1
   
`Run the aws command to ssh`
    
        docker run -it -v $(dirname $SSH_AUTH_SOCK):$(dirname $SSH_AUTH_SOCK) -e SSH_AUTH_SOCK=$SSH_AUTH_SOCK  -v ~/.ssh:/root/.ssh/ -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} -e AWS_REGION=${AWS_REGION} --entrypoint sh mwaaas/aws_ssh:latest
        
     commit in tag 
