运行前,需要先将需要运行的容器镜像cp至`RootUrl`目录
例如busybox
```shell
docker run -d --name busybox-run1 busybox top -b
docker export -o busybox.tar busybox-run1 && mkdir -p busybox && tar -xvf busybox.tar -C busybox/ 
```
