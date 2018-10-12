package types

type ModelConfig struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	ModelPath   string `json:"model_path" yaml:"model_path"`
}
