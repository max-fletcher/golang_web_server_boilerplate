package common_errors

// NOTE: rule of thumb for defining custom errors:
// 1. Define error interfaces for dynamic errors(errors that need to be constructed using variables)(see below example)
// 2. Define static errors for static strings(see below example)
// rule of thumb for handling custom errors:
// 1. use "if err != nil{...}" for static one-off errors(error that contain only string and you don't intend to return it
//    conditionally e.g you don't want to determine which status code should be thrown)
// 2. use "if errors.Is(err, ErrUserNotFound) {...}" for error static errors that you want to return conditionally(i.e you
//    want to determine what status code should be returned with response from handler/controller)
// 3. use
// "if err != nil {
// 	var emailExistsErr users.ErrUserWithEmailAlreadyExists
// 	if errors.As(err, &emailExistsErr) {
// 		responses.ConflictError(w, emailExistsErr.Error())
// 		return
// 	}
//  ...
// }"
// for dynamic errors and you want to return response conditionally(i.e determine which status code to throw based on what
// error type). errors.As stores the err into &emailExistsErr if the type matches and from there, you can throw it how you please

type ValidationError struct {
	Errors map[string]string
}

func (e ValidationError) Error() string {
	return "Validation failed"
}

type ErrHashingPassword struct {
	HashErr error
}

func (e ErrHashingPassword) Error() string {
	return "Failed to hash password"
}

func (e ErrHashingPassword) Unwrap() error {
	return e.HashErr
}

var (
// common
// ErrNotLoggedIn            = errors.New("You are not logged in")
// ErrUnauthorized           = errors.New("You are unauthorized to access this resource")
// ErrInsufficientPermission = errors.New("You don't have sufficient permission to access this resource")
// ErrHashingPassword = errors.New("Error processing password")
)
