export type WithChildren<T = unknown> = {
    children?: React.ReactElement | React.ReactElement[] | string
} & T
