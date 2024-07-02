echo "Please Run Me on the root dir, not in scripts dir."

if [ -d "output" ]; then
    echo "Output dir existed, deleting and recreating..."
    rm -rf output
fi
mkdir -p output/services

pushd backend/go-services || exit

for i in *; do
  if [ -d "$i" ] && [ "$i" != "health" ]; then
      echo "$i"
      name="$i"
      capName="${name}"
      cd "$i"/services
      go build -o "../../../../output/services/$i/${capName}Service"
      cd ../..
  fi
done

popd || exit

mkdir -p output/gateway

cd backend/api-gateway || exit

go build -o "../../output/gateway/Gateway"

echo "OK!"