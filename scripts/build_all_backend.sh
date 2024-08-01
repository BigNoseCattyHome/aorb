# ! /bin/bash

echo "Please Run Me on the root dir, not in scripts dir."

AORB_HOME=$(pwd)

# 清理 build 目录
if [ -d "build" ]; then
  echo "./build/ dir existed, removing and recreating..."
  rm -rf build
fi

# 构建 gateway
#cd $AORB_HOME
#mkdir -p build/gateway
#cd backend/api-gateway || exit
#go build -o "$AORB_HOME/build/gateway/Gateway"

# 构建 services
cd $AORB_HOME
mkdir -p build/services
cd backend/go-services || exit

for i in *; do
  if [ -d "$i" ]; then
    name="$i"
    capName="${name}"
    cd "$i"
    echo "Building $i service..."
    go build -o "$AORB_HOME/build/services/$i/${capName}Service"
    cd ..
  fi
done

echo "Build Done!"