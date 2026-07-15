plugins {
    kotlin("jvm") version "2.2.21"
    `java-library`
    `maven-publish`
}

group = "com.auth0.kotlin"
version = "1.0.0"

dependencies {
    api("com.auth0:auth0:3.10.0")
    api("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.10.2")
    testImplementation(kotlin("test"))
    testImplementation("org.junit.jupiter:junit-jupiter:6.0.3")
}

val integrationTest by sourceSets.creating

integrationTest.compileClasspath += sourceSets.main.get().output
integrationTest.runtimeClasspath += sourceSets.main.get().output

configurations[integrationTest.implementationConfigurationName].extendsFrom(configurations.testImplementation.get())
configurations[integrationTest.runtimeOnlyConfigurationName].extendsFrom(configurations.testRuntimeOnly.get())

kotlin {
    explicitApi()
    jvmToolchain(17)
    compilerOptions {
        freeCompilerArgs.add("-Xjsr305=strict")
    }
}

java {
    withSourcesJar()
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
}

tasks.test {
    useJUnitPlatform()
}

tasks.register<Test>("integrationTest") {
    testClassesDirs = integrationTest.output.classesDirs
    classpath = integrationTest.runtimeClasspath
    useJUnitPlatform()
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            from(components["java"])
            artifactId = "auth0-kotlin"
        }
    }
}
