package internal

import libUUID "github.com/google/uuid"

const DOCKER_CONTAINER_ID_FIXED_LENGTH int = 64
const DOCKER_CONTAINER_ID_TRUNCATED int = 12

const EPHEMERAL_IMAGE_NAME = "a-mule-ephemeral"
const ENGINE_IMAGE_NAME = "a-mule-engine"

var UUID_NAMESPACE = libUUID.MustParse("b6ace449-d6e4-4dd6-a87b-821a1523a374")
