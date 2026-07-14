@echo off
set APP_HOME=%~dp0
set JAVACMD=java
if not "%JAVA_HOME%" == "" set JAVACMD=%JAVA_HOME%\bin\java
"%JAVACMD%" -classpath "%APP_HOME%\gradle\wrapper\gradle-wrapper.jar" org.gradle.wrapper.GradleWrapperMain %*
