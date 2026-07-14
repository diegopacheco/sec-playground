pluginManagement {
    repositories {
        gradlePluginPortal()
        mavenCentral()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        mavenCentral()
        maven("https://s01.oss.sonatype.org/content/repositories/releases/")
    }
}

rootProject.name = "auth0-kotlin-2x"

includeBuild("/Users/diegopacheco/git/misc/auth0-java") {
    dependencySubstitution {
        substitute(module("com.auth0:auth0")).using(project(":"))
    }
}
