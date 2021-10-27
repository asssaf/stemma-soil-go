: ${IMAGE_NAME:=asssaf/stemma-soil:latest}
docker run --rm -it --privileged --device /dev/i2c-1 "$IMAGE_NAME" $*
