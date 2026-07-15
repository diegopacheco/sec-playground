plugins {
    kotlin("jvm")
    application
}

dependencies {
    implementation(rootProject)
}

kotlin {
    jvmToolchain(17)
}

application {
    mainClass.set("com.auth0.kotlin.webapp.ServerKt")
}

tasks.named<JavaExec>("run") {
    workingDir = File(rootProject.projectDir, "sample-app")
}
