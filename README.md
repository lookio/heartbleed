heartbleed
==========

To build the server, use Docker and run

```
sudo docker build .
```

Then run the server with:

```
sudo docker run -t -i -p 8000:8000 -v `pwd`:/code <image_id> /bin/bash
```

Inside the Docker container:

```
cd /code
go build .
./code -v
```

This is for running a dev build, running in prod is different!
