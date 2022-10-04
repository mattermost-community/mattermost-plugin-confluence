export type ConfluenceConfig = {
    serverURL: string,
    clientID: string,
    clientSecret: string
}

export type ReactSelectOption = {
    label: string | React.ReactElement;
    value: string;
};

export type ErrorResponse = {
    response: {
        text: string;
    }
};
