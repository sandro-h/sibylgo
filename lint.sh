golint ./... | grep -v vendor > lint_report.txt
if [ -s lint_report.txt ]; then
    echo "Lint failed"
    exit 1
fi
echo "Lint OK"