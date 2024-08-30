/*
 * This file was generated by the Gradle 'init' task.
 *
 * This is a general purpose Gradle build.
 * To learn more about Gradle by exploring our Samples at https://docs.gradle.org/8.7/samples
 */
plugins {
    id("com.fussionlabs.gradle.go-plugin") version("0.7.0")
}

version = "0.0.1"
val moduleName = "github.com/dm0275/mcrun"

go {
    os = listOf("linux", "darwin")
    arch = listOf("amd64", "arm64")
    ldFlags = mapOf("github.com/dm0275/mcrun/pkg/version.version" to "$version")
}