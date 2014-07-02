function run_tests {
# maxreport
##
  tool="./maxreport/maxreport.go"

  output=`go run $tool -h`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  output=`go run $tool`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  if test "$ALIAS" && test "$TOKEN" && test "$SECRET"
  then # Run functional tests.
    output=`go run $tool stats`
    assert_grep "echo '$output'" "Running summary stats report." \
      "$tool: returns expected header"
    assert_grep "echo '$output'" "total hits" \
      "$tool: returns expected output"

    output=`go run $tool popular`
    assert_grep "echo '$output'" "Running popular files report." \
      "$tool: returns expected header"
    assert_grep "echo '$output'" "hits | file" \
      "$tool: returns expected output"

    assert "go run $tool popular --top 5 | grep '[0-9]+ \| \/' | wc -l" \
      "$tool: returns expected number of results"
  fi
}

