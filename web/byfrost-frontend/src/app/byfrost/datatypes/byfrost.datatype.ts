
// Given this rubbish name to accommodate the fact that this might sometimes be
// a class and not just functions.
export interface ProgramConstruct {
    name: string;
    file: string;
    tokens: Token[];
}

interface Token {
    name: string;
    filepath: string;
    typ: string;
    spaces: number;
    tabs: number;
    meta: string;
}