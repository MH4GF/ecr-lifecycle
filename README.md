# ecr-lifecycle

- ECRのイメージを,指定した件数より古いバージョンを削除
- ECSタスクとして現在実行されている場合、削除せず保護

# Usage

```shell script
$ ecr-lifecycle delete-images --keep 50 --ecr-profile sandbox --region ap-northeast-1 --ecs-profiles hoge,fuga
```

# development

```shell script
$ make build
$ bin/ecr-lifecycle
```
