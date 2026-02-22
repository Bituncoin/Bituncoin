plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("plugin.serialization") version "1.9.10"
    application
}

group = "global.genesis"
version = "1.0.0"
java.sourceCompatibility = JavaVersion.VERSION_17

application {
    mainClass.set("global.genesis.BtngGenesisPlatformKt")
}

repositories {
    mavenCentral()
}

dependencies {
    // Kotlin Standard Library
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")

    // JSON processing
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.0")

    // HTTP Client for BTNG API integration
    implementation("io.ktor:ktor-client-core:2.3.0")
    implementation("io.ktor:ktor-client-cio:2.3.0")
    implementation("io.ktor:ktor-client-content-negotiation:2.3.0")
    implementation("io.ktor:ktor-serialization-kotlinx-json:2.3.0")

    // WebSocket support for real-time streaming
    implementation("io.ktor:ktor-client-websockets:2.3.0")

    // Database drivers (H2 for development)
    implementation("com.h2database:h2:2.1.214")

    // Logging
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("org.slf4j:slf4j-api:2.0.9")

    // Coroutines for async operations
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3")

    // Testing
    testImplementation("org.junit.jupiter:junit-jupiter:5.10.0")
    testImplementation("io.kotest:kotest-runner-junit5:5.6.2")
    testImplementation("com.ninja-squad:springmockk:4.0.2")
}

tasks.test {
    useJUnitPlatform()
}

tasks.jar {
    manifest {
        attributes["Main-Class"] = "global.genesis.BtngGenesisPlatformKt"
    }
    duplicatesStrategy = DuplicatesStrategy.EXCLUDE
    from(sourceSets.main.get().output)
    dependsOn(configurations.runtimeClasspath)
    from({
        configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) }
    })
}

tasks.test {
    useJUnitPlatform()
}

tasks.wrapper {
    gradleVersion = "8.5"
    distributionType = Wrapper.DistributionType.ALL
}

// Genesis-specific tasks
tasks.register("genesisRun") {
    dependsOn("build")
    doLast {
        exec {
            workingDir = projectDir
            commandLine("java", "-cp", "build/libs/${project.name}-${project.version}.jar", "global.genesis.BtngGenesisPlatformKt")
        }
    }
}