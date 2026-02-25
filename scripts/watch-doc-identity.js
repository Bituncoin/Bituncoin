#!/usr/bin/env node

const { watch } = require('fs')
const { resolve } = require('path')
const { generateDocumentationIdentity } = require('./generate-doc-identity')
const { propagateDocumentationIdentity } = require('./propagate-doc-identity')

const ROOT_DIR = process.cwd()
const CONFIG_PATH = resolve(ROOT_DIR, 'docs/canonical-docs.json')

function loadCanonicalFiles() {
    const config = require(CONFIG_PATH)
    return Array.isArray(config.files) ? config.files : []
}

let timer = null

function regenerate() {
    try {
        generateDocumentationIdentity()
        propagateDocumentationIdentity()
        console.log('DOC_IDENTITY_WATCH=UPDATED')
    } catch (error) {
        console.error(`DOC_IDENTITY_WATCH=FAIL\n${error.message}`)
    }
}

function scheduleRegeneration() {
    if (timer) {
        clearTimeout(timer)
    }
    timer = setTimeout(regenerate, 300)
}

function startWatcher() {
    const files = loadCanonicalFiles()
    if (files.length === 0) {
        throw new Error('No canonical documentation files configured.')
    }

    regenerate()

    files.forEach((relativePath) => {
        const absolutePath = resolve(ROOT_DIR, relativePath)
        watch(absolutePath, { persistent: true }, () => scheduleRegeneration())
    })

    watch(CONFIG_PATH, { persistent: true }, () => {
        console.log('Canonical docs configuration changed. Restart watcher to include new files.')
        scheduleRegeneration()
    })

    console.log(`DOC_IDENTITY_WATCH=ACTIVE files=${files.length}`)
}

if (require.main === module) {
    try {
        startWatcher()
    } catch (error) {
        console.error(`DOC_IDENTITY_WATCH=FAIL\n${error.message}`)
        process.exit(1)
    }
}

module.exports = {
    startWatcher
}
