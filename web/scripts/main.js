"use strict";

export function getRegisterData() {
  return {
    username: "",
    password: "",
    confirmPassword: "",
    pwLen: false,
    pwEvalLen: false,
    pwMixedCase: false,
    pwEvalMixedCase: false,
    pwNum: false,
    pwEvalNum: false,
    pwSpecialChar: false,
    pwEvalSpecialChar: false,
    submitData,
    sanitizeUsername,
    isUsernameValid,
    isPasswordValid,
    getPasswordStrengths,
  };
}

export async function submitData(u, pw, confirmpw) {
  if (!isUsernameValid(u)) return;
  if (!isPasswordValid(pw)) return;
  if (pw !== confirmpw) return;
  const su = sanitizeUsername(u);

  const body = new FormData();
  body.append("username", su);
  body.append("password", pw);

  try {
    const response = await fetch("/player/new", {
      method: "POST",
      body,
    });

    if (response.status !== 201) {
      // TODO: Handle 409 here
      // TODO: Handle 500 and 400 as server-side errors here
      return;
    }

    window.location.reload();
  } catch (err) {
    // TODO: Here, the fetch has failed - this is a network partition or the back end went down
    return;
  }
}

export function sanitizeUsername(u) {
  return u.replace(/[^a-zA-Z0-9_\-]+/gi, "").toLowerCase();
}

export function isUsernameValid(u) {
  if (u.length < 4) return false;
  if (u.length > 16) return false;
  const regex = new RegExp("[^a-z0-9_-]+", "g");
  if (regex.test(u)) return false;
  return true;
}

export function isPasswordValid(pw) {
  if (pw.length < 8) return false;
  if (pw.length > 255) return false;
  return true;
}

export function getPasswordStrengths(pw) {
  let strengths = {
    len: false,
    mixedCase: false,
    num: false,
    specialChar: false,
  };

  if (pw.length > 8) {
    strengths.len = true;
  }

  if (pw.match(/[a-z]/) && pw.match(/[A-Z]/)) {
    strengths.mixedCase = true;
  }

  if (pw.match(/[0-9]/)) {
    strengths.num = true;
  }

  if (pw.match(/[^a-zA-Z\d]/)) {
    strengths.specialChar = true;
  }

  return strengths;
}

window.getRegisterData = getRegisterData;
