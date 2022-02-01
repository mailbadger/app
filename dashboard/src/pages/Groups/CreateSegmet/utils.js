import { string, object } from "yup"

export const segmentValidation = object().shape({
    name: string()
        .required("Please enter a group name.")
        .max(191, "The name must not exceed 191 characters."),
})
