"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
const SignupButton = document.getElementById("signUpButton");
const SignupField = document.getElementById("signUpField");
const ErrorMessage = document.getElementById("errorMessage");
const ErrorMessageText = document.getElementById("errorMessageText");
const ERROR_LENGTH_LESS = 0;
const ERROR_LENGTH_MORE = 1;
const ERROR_INVALID_CHAR = 2;
const ERROR_UDSC_PRD_MORE = 3;
const ERROR_IN_USE = 4;
const ERROR_UNEXPECTED = 5;
SignupButton.addEventListener("click", (e) => {
    const username = SignupField.value;
    EvaluateUsername(username);
});
function EvaluateUsername(usnm) {
    SignupField.classList.remove("hasError");
    ErrorMessageText.innerText = "";
    ErrorMessage.classList.add("hide");
    // LENGTH RELATED
    if (usnm.length < 4)
        return ShowError(ERROR_LENGTH_LESS);
    if (usnm.length > 14)
        return ShowError(ERROR_LENGTH_MORE);
    // INVALID CHAR
    if (!/^[a-z0-9._]+$/i.test(usnm))
        return ShowError(ERROR_INVALID_CHAR);
    // UNDERSCORE & PERIOD
    const regMatch = usnm.match(/[._]/g);
    if (regMatch && regMatch.length > 1)
        return ShowError(ERROR_UDSC_PRD_MORE);
    SetUsername(usnm);
}
function ShowError(errorID) {
    SignupField.classList.add("hasError");
    ErrorMessage.classList.remove("hide");
    console.log("Show error:", errorID);
    let errMsg = "";
    switch (errorID) {
        case ERROR_LENGTH_LESS:
            errMsg = "Your username must be at least 4 characters long.";
            break;
        case ERROR_LENGTH_MORE:
            errMsg = "Your username must not exceed 14 characters.";
            break;
        case ERROR_INVALID_CHAR:
            errMsg = "Your username contains invalid characters. It can only contain letters, numbers, and one underscore or period.";
            break;
        case ERROR_UDSC_PRD_MORE:
            errMsg = "Your username cannot contain more than one underscore or one period.";
            break;
        case ERROR_IN_USE:
            errMsg = "That username is in use. Please choose a different one.";
            break;
        case ERROR_UNEXPECTED:
            errMsg = "There has been an internal error. Please try again later.";
            break;
    }
    ErrorMessageText.innerText = errMsg;
}
function SetUsername(usnm) {
    return __awaiter(this, void 0, void 0, function* () {
        let response = yield fetch(`/api/endSignup`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username: usnm }),
        });
        if (response.ok) {
            if (response.status == 200) {
                window.location.href = "/";
            }
        }
        else {
            let result = yield response.json();
            if (result.errorID !== undefined) {
                ShowError(result.errorID);
            }
            else {
                ShowError(ERROR_UNEXPECTED);
            }
        }
    });
}
