enum ExprType {
    Expr = "Expr",
    Eq = "Eq"
}

type MathReq = {
    exprType: ExprType,
    expr: string
}

type MathRes = {
    exprType: ExprType,
    expr: Array<string>
}

export {
    MathReq,
    MathRes,
    ExprType
}