mkdir docs
dir=$(pwd)
for d in */; do
    cd $dir
    echo $d
    cd $d
    serviceName=${d//\//}
    timeout 3s make proto || continue
    echo "Copying html for $serviceName"
    cp redoc-static.html ../docs/$serviceName-api.html || continue
    echo "---\ntitle: $servicename\n---\n"../docs/hugo-tania/exampleSite/content/post/$serviceName.md || continue
    cat README.md > ../docs/hugo-tania/exampleSite/content/post/$serviceName.md || continue
done
pwd
cd ../docs/hugo-tania/exampleSite; hugo -D -d=../../
cd ../../
pwd
ls