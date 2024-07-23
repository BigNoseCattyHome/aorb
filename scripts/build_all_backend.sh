echo "Please Run Me on the root dir, not in scripts dir."

if [ -d "build" ]; then
    echo "Output dir existed, deleting and recreating..."
    rm -rf build
fi
mkdir -p build/services

pushd backend/go-services || exit

for i in *; do
  if [ -d "$i" ] && [ "$i" != "health" ]; then
      echo "$i"
      name="$i"
      capName="${name}"
      cd "$i"/services
      go build -o "../../../../build/services/$i/${capName}Service"
      cd ../..
  fi
done

popd || exit

mkdir -p build/gateway

cd backend/api-gateway || exit

go build -o "../../build/gateway/Gateway"

echo "OK!"