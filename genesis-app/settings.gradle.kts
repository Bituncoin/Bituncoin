rootProject.name = "btng-genesis-app"

pluginManagement {
    repositories {
        mavenCentral()
        gradlePluginPortal()
        maven {
            url = uri("https://genesisglobal.jfrog.io/artifactory/gg-maven-public")
            credentials {
                username = System.getenv("ARTIFACTORY_USER") ?: ""
                password = System.getenv("ARTIFACTORY_PASSWORD") ?: ""
            }
        }
    }
}

dependencyResolutionManagement {
    repositories {
        mavenCentral()
        maven {
            url = uri("https://genesisglobal.jfrog.io/artifactory/gg-maven-public")
            credentials {
                username = System.getenv("ARTIFACTORY_USER") ?: ""
                password = System.getenv("ARTIFACTORY_PASSWORD") ?: ""
            }
        }
    }
}

// Include subprojects if they exist
include("server")
include("client")