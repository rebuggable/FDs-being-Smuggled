# FDs being smuggled into a container instance

It is supposed to enable you to mount arbitrary directories into a running container instance without the need to recreate or restart the container.

## Demonstration by Example

Alex would love to share a folder at `C:\Users\alex\Documents\dataset` with an agent inside a container without restarting or recreating the container. The agent expects the files to be available at `/mnt/curated/dataset`

```
mule 99799fd4589a --mount C:\Users\alex\Documents\dataset:/mnt/curated/dataset
```

> You can pass as many `--mount` arguments as you'd like.

In this example, `/mnt/curated/dataset` is the mount point, and it must exist before you execute the command. If it doesn't exist, you run
```
mkdir -p /mnt/curated/dataset
```

We actually do it for you if the directory doesn't exist, but there is a caveat here. Once a container restarts, the directories will be unmounted, but the mountpoint won't.

### Unmount directories

You can unmount directories using the `--unmount` flag. All you need to do is assign it the directory value inside the container instance, like so

```
mule 99799fd4589a --unmount /mnt/curated/dataset
```

## Docker in Docker

Yeah, I haven't finished it yet. Neither have I finished the other one, LOL. But the simple reason is that I'm forced to set the PID mode to Host. I didn't like exposing the other PIDs. So, if you don't like this either, you can run the whole situation in Docker-in-Docker mode and run the command as follows

```
mule a7f3c91e4b82/99799fd4589a --mount /mnt/Doc/dataset:/mnt/curated/dataset
```

While `a7f3c91e4b82` is the DinD container ID, `99799fd4589a` is the inner container instance ID (inside the DinD container with ID `a7f3c91e4b82`).

Finally, you're not forced to provide the container ID or name as a CLI argument; you can pipe it in :D

```
... | mule --mount ...
```

> Please use the DinD image under `Dockerfiles/engine` because it already has the orchestrator deployed and configured for use.