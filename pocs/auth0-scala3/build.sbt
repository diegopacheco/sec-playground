ThisBuild / scalaVersion := "3.7.2"
ThisBuild / organization := "com.auth0.scala3"
ThisBuild / version := "1.0.0"

lazy val root = (project in file("."))
  .settings(
    name := "auth0-scala-3x",
    scalacOptions ++= Seq("-deprecation", "-feature"),
    libraryDependencies += "com.auth0" % "auth0" % "3.10.0",
    Test / test := {
      (Test / compile).value
      val cp = (Test / fullClasspath).value.files.map(_.getAbsolutePath).mkString(java.io.File.pathSeparator)
      val exit = sys.process.Process(Seq("java", "-cp", cp, "auth0scala3.Auth0Spec")).!
      if (exit != 0) sys.error("tests failed")
    }
  )
