ThisBuild / scalaVersion := "3.7.2"
ThisBuild / organization := "com.auth0.scala3"
ThisBuild / version := "1.0.0"

lazy val IntegrationTest = config("it") extend Test

lazy val root = (project in file("."))
  .configs(IntegrationTest)
  .settings(
    name := "auth0-scala-3x",
    scalacOptions ++= Seq("-deprecation", "-feature", "-unchecked", "-Wunused:all"),
    libraryDependencies ++= Seq(
      "com.auth0" % "auth0" % "3.10.0",
      "org.scalameta" %% "munit" % "1.2.0" % Test
    )
  )
  .settings(inConfig(IntegrationTest)(Defaults.testSettings))
