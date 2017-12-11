FROM scratch
ADD /env2file /bin/env2file
ENTRYPOINT ["env2file"]