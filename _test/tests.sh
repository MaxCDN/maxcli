function run_tests {

  for tool in `find . -maxdepth 2 -type f -name "max*.go"`
  do
    assert_grep "go run $tool -h" "Usage" \
      "deplay usage without params"
  done

}
