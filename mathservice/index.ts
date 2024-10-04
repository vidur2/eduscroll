import { parse } from "path";
import { ExprType, MathReq, MathRes } from "./types/mathreqs";

const express = require("express");
const mathsteps = require("mathsteps");
const app = express();
const port = 8080;

app.use(express.json())

app.get("/", (req: any, res: any) => {
    res.status(400);
    res.json({
        hello: "world"
    });
})

app.post("/", (req: any, res: any) => {
    const body: MathReq = req.body
    let out: Array<string>;
    switch (body.exprType) {
        case ExprType.Eq: 
            out = mathsteps.solveEquation(body.expr).map((step) => { return { old: step.oldEquation.latex().toString(), new: step.newEquation.latex().toString() } });
            break;
        case ExprType.Expr:
            out = mathsteps.simplifyExpression(body.expr).map((step) => { return { old: step.oldNode.toString(), new: step.newNode.toString() } });
            break
    }

    const outJson: MathRes = {
        exprType: body.exprType,
        expr: out
    }
    res.json(outJson)
})

app.listen(port, "0.0.0.0", () => {
    console.log(`listening on ${port}`);
})