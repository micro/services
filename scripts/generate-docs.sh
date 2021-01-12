mkdir docs
dir=$(pwd)
for d in */; do
    cd $dir
    echo $d
    cd $d
    timeout 3s make proto || continue
    echo "Copying html for $d"
    cp redoc-static.html ../docs/${d//\//}-api.html || continue
    cp README.md ../docs/hugo-tania/exampleSite/content/post/${d//\//}.md || continue
done
ls docs