package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "testing"
)

func TestTruncateLinesWidth(t *testing.T) {
    str := "12345678\nABCDEFGH\nX\n\n"
    assertEquals(t, truncateLinesWidth(str, 7), "1234...\nABCD...\nX\n\n")
    assertEquals(t, truncateLinesWidth(str, 6), "123...\nABC...\nX\n\n")
    assertEquals(t, truncateLinesWidth(str, 5), "12...\nAB...\nX\n\n")
    assertEquals(t, truncateLinesWidth(str, 4), "1...\nA...\nX\n\n")
}

func TestTruncateLinesCount(t *testing.T) {
    str := "1\n2\n3\n4\n5\n6\n7\n8\n9"
    assertEquals(t, truncateLinesCount(str, 10), "1\n2\n3\n4\n5\n6\n7\n8\n9")
    assertEquals(t, truncateLinesCount(str, 9), "1\n2\n3\n4\n5\n6\n7\n8\n9")
    assertEquals(t, truncateLinesCount(str, 8), "1\n2\n3\n4\n5\n6\n7\n8\n...")
    assertEquals(t, truncateLinesCount(str, 4), "1\n2\n3\n4\n...")
}

func TestMatchEdgeCases(t *testing.T) {
    assertNonMatch(t, "")
    assertNonMatch(t, "\n")
    assertNonMatch(t, "abc")
}

func TestMatchJavaCompilationError(t *testing.T) {
    start := "[ERROR] COMPILATION ERROR"
    end := "[INFO] ------------------------------------------------------------------------"
    bodyStr := start + "\n" + "bla bla bla" + "\n" + end
    assertMatch(t, bodyStr, "Bambot detected a Java compilation error!")
}

func TestCSharpError(t *testing.T) {
    fileName := "test_files/csharp-1.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected a C# build error!")
}

func TestPythonError(t *testing.T) {
    fileName := "test_files/python-1.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected a Python error!")
}

func TestGenericError(t *testing.T) {
    fileName := "test_files/generic.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected an error!")
}

// When there are multiple matches, we want to identify only the last match present in the log file
func TestMultipleMatchesGenericError(t *testing.T) {
    fileName := "test_files/generic-multiple-matches.log"
    bodyStr := readFileToString(fileName)
    scanResult := assertMatch(t, bodyStr, "Bambot detected an error!")
    assertContains(t, scanResult.LogSnippet, "<this should be included>")
    assertNotContains(t, scanResult.LogSnippet, "<this should not be included>")
}

func TestGruntError(t *testing.T) {
    fileName := "test_files/grunt-1.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected a front-end Grunt build error!")
}

func TestCSharpTestError(t *testing.T) {
    fileName := "test_files/csharp-test-1.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected a C# unit test/integration test failure!")
}

func TestPythonTestError(t *testing.T) {
    fileName := "test_files/python-test-1.log"
    bodyStr := readFileToString(fileName)
    assertMatch(t, bodyStr, "Bambot detected a Python (sdk?) unit test test failure!")
}

func assertEquals(t *testing.T, str string, expectedStr string) string {
    if str != expectedStr {
        t.Errorf("expected '%s' but got '%s'", str, expectedStr)
    }
    return str
}

func assertContains(t *testing.T, str string, subStr string) string {
    if strings.Index(str, subStr) < 0 {
        t.Errorf("expected '%s' to contain '%s' but it did not", str, subStr)
    }
    return str
}

func assertNotContains(t *testing.T, str string, subStr string) string {
    if strings.Index(str, subStr) >= 0 {
        t.Errorf("expected '%s' to not contain '%s' but it did", str, subStr)
    }
    return str
}

func assertNonMatch(t *testing.T, bodyStr string) ScanResult {
    scanResult := scanString(bodyStr)
    if scanResult != nonMatch() {
       t.Errorf("expected '%s' to not match any rules, result was '%s'", bodyStr, scanResult)
    }
    return scanResult
}

func debugPrintScanResult(scanResult ScanResult) {
    fmt.Println("\tComment: " + scanResult.Comment)
    fmt.Println("\tJIRA:" + scanResult.JiraIssueId)
    fmt.Println("\tSnippet:\n\t\t" + strings.Join(strings.Split(scanResult.LogSnippet, "\n"), "\n\t\t"))
}

func assertMatch(t *testing.T, bodyStr string, expectedComment string) ScanResult {
    scanResult := scanString(bodyStr)
    if scanResult == nonMatch() {
        t.Errorf("expected '%s' to match a rule, but it matched nothing", truncate(bodyStr))
    }
    if scanResult.Comment != expectedComment {
        t.Errorf("expected comment '%s' but found '%s'", expectedComment, scanResult.Comment)
    }

    return scanResult
}

func truncate(bodyStr string) string {
    if len(bodyStr) > 50 {
        return bodyStr[0:50]
    } else {
        return bodyStr
    }
}

func readFileToString(fileName string) string {
    content, err := ioutil.ReadFile(fileName)
    if err != nil {
        panic(err)
    }
    return string(content)
}