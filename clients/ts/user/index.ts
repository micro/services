import * as m3o from "@m3o/m3o-node";

export class UserService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Create a new user account. The email address and username for the account must be unique.
  create(request: CreateRequest): Promise<CreateResponse> {
    return this.client.call(
      "user",
      "Create",
      request
    ) as Promise<CreateResponse>;
  }
  // Delete an account by id
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call(
      "user",
      "Delete",
      request
    ) as Promise<DeleteResponse>;
  }
  // Login using username or email. The response will return a new session for successful login,
  // 401 in the case of login failure and 500 for any other error
  login(request: LoginRequest): Promise<LoginResponse> {
    return this.client.call("user", "Login", request) as Promise<LoginResponse>;
  }
  // Logout a user account
  logout(request: LogoutRequest): Promise<LogoutResponse> {
    return this.client.call(
      "user",
      "Logout",
      request
    ) as Promise<LogoutResponse>;
  }
  // Read an account by id, username or email. Only one need to be specified.
  read(request: ReadRequest): Promise<ReadResponse> {
    return this.client.call("user", "Read", request) as Promise<ReadResponse>;
  }
  // Read a session by the session id. In the event it has expired or is not found and error is returned.
  readSession(request: ReadSessionRequest): Promise<ReadSessionResponse> {
    return this.client.call(
      "user",
      "ReadSession",
      request
    ) as Promise<ReadSessionResponse>;
  }
  // Send a verification email
  // to the user being signed up. Email from will be from 'support@m3o.com',
  // but you can provide the title and contents.
  // The verification link will be injected in to the email as a template variable, $micro_verification_link.
  // Example: 'Hi there, welcome onboard! Use the link below to verify your email: $micro_verification_link'
  // The variable will be replaced with an actual url that will look similar to this:
  // 'https://user.m3o.com/user/verify?token=a-verification-token&redirectUrl=your-redir-url'
  sendVerificationEmail(
    request: SendVerificationEmailRequest
  ): Promise<SendVerificationEmailResponse> {
    return this.client.call(
      "user",
      "SendVerificationEmail",
      request
    ) as Promise<SendVerificationEmailResponse>;
  }
  // Update the account password
  updatePassword(
    request: UpdatePasswordRequest
  ): Promise<UpdatePasswordResponse> {
    return this.client.call(
      "user",
      "UpdatePassword",
      request
    ) as Promise<UpdatePasswordResponse>;
  }
  // Update the account username or email
  update(request: UpdateRequest): Promise<UpdateResponse> {
    return this.client.call(
      "user",
      "Update",
      request
    ) as Promise<UpdateResponse>;
  }
  // Verify the email address of an account from a token sent in an email to the user.
  verifyEmail(request: VerifyEmailRequest): Promise<VerifyEmailResponse> {
    return this.client.call(
      "user",
      "VerifyEmail",
      request
    ) as Promise<VerifyEmailResponse>;
  }
}

export interface Account {
  // unix timestamp
  created?: number;
  // an email address
  email?: string;
  // unique account id
  id?: string;
  // Store any custom data you want about your users in this fields.
  profile?: { [key: string]: string };
  // unix timestamp
  updated?: number;
  // alphanumeric username
  username?: string;
  verificationDate?: number;
  verified?: boolean;
}

export interface CreateRequest {
  // the email address
  email?: string;
  // optional account id
  id?: string;
  // the user password
  password?: string;
  // optional user profile as map<string,string>
  profile?: { [key: string]: string };
  // the username
  username?: string;
}

export interface CreateResponse {
  account?: { [key: string]: any };
}

export interface DeleteRequest {
  // the account id
  id?: string;
}

export interface DeleteResponse {}

export interface LoginRequest {
  // The email address of the user
  email?: string;
  // The password of the user
  password?: string;
  // The username of the user
  username?: string;
}

export interface LoginResponse {
  // The session of the logged in  user
  session?: { [key: string]: any };
}

export interface LogoutRequest {
  sessionId?: string;
}

export interface LogoutResponse {}

export interface ReadRequest {
  // the account email
  email?: string;
  // the account id
  id?: string;
  // the account username
  username?: string;
}

export interface ReadResponse {
  account?: { [key: string]: any };
}

export interface ReadSessionRequest {
  // The unique session id
  sessionId?: string;
}

export interface ReadSessionResponse {
  session?: { [key: string]: any };
}

export interface SendVerificationEmailRequest {
  email?: string;
  failureRedirectUrl?: string;
  // Display name of the sender for the email. Note: the email address will still be 'support@m3o.com'
  fromName?: string;
  redirectUrl?: string;
  subject?: string;
  // Text content of the email. Don't forget to include the string '$micro_verification_link' which will be replaced by the real verification link
  // HTML emails are not available currently.
  textContent?: string;
}

export interface SendVerificationEmailResponse {}

export interface Session {
  // unix timestamp
  created?: number;
  // unix timestamp
  expires?: number;
  // the session id
  id?: string;
  // the associated user id
  userId?: string;
}

export interface UpdatePasswordRequest {
  // confirm new password
  confirmPassword?: string;
  // the new password
  newPassword?: string;
  // the old password
  oldPassword?: string;
  // the account id
  userId?: string;
}

export interface UpdatePasswordResponse {}

export interface UpdateRequest {
  // the new email address
  email?: string;
  // the account id
  id?: string;
  // the user profile as map<string,string>
  profile?: { [key: string]: string };
  // the new username
  username?: string;
}

export interface UpdateResponse {}

export interface VerifyEmailRequest {
  // The token from the verification email
  token?: string;
}

export interface VerifyEmailResponse {}
