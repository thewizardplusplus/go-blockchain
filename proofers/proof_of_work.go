package proofers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/mo"
	"github.com/thewizardplusplus/go-blockchain"
	pow "github.com/thewizardplusplus/go-pow"
	powErrors "github.com/thewizardplusplus/go-pow/errors"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

const (
	hashPartSeparator = ":"
	hashPartCount     = 3
	maximalTargetBit  = sha256.Size*8 - 1
)

var ErrInvalidParameters = errors.New("invalid parameters")

// ProofOfWork ...
type ProofOfWork struct {
	TargetBit                int
	MaxAttemptCount          mo.Option[int]
	RandomInitialNonceParams mo.Option[powValueTypes.RandomNonceParams]
}

// Hash ...
//
// Deprecated: Use [ProofOfWork.HashEx] instead.
func (proofer ProofOfWork) Hash(block blockchain.Block) string {
	hash, _ := proofer.HashEx(context.Background(), block)
	return hash
}

// HashEx ...
func (proofer ProofOfWork) HashEx(
	ctx context.Context,
	block blockchain.Block,
) (string, error) {
	targetBitIndex, err := powValueTypes.NewTargetBitIndex(proofer.TargetBit)
	if err != nil {
		return "", fmt.Errorf(
			"unable to construct the target bit index: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	challenge, err := pow.NewChallengeBuilder().
		SetTargetBitIndex(targetBitIndex).
		SetSerializedPayload(powValueTypes.NewSerializedPayload(block.MergedData())).
		SetHash(powValueTypes.NewHash(sha256.New())).
		SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
			"{{ .Challenge.SerializedPayload.ToString }}" +
				"{{ .Nonce.ToString }}" +
				"{{ .Challenge.TargetBitIndex.ToInt }}",
		)).
		Build()
	if err != nil {
		return "", fmt.Errorf(
			"unable to build the challenge: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	solution, err := challenge.Solve(ctx, pow.SolveParams{
		MaxAttemptCount:          proofer.MaxAttemptCount,
		RandomInitialNonceParams: proofer.RandomInitialNonceParams,
	})
	if err != nil {
		if !errors.Is(err, powErrors.ErrIO) &&
			!errors.Is(err, powErrors.ErrTaskInterruption) {
			err = errors.Join(err, ErrInvalidParameters)
		}

		return "", fmt.Errorf("unable to solve the challenge: %w", err)
	}

	hashSum, isPresent := solution.HashSum().Get()
	if !isPresent {
		return "", fmt.Errorf("hash sum is absent in the solution: %w", err)
	}

	hashParts := []string{
		strconv.Itoa(proofer.TargetBit),
		solution.Nonce().ToString(),
		hex.EncodeToString(hashSum.ToBytes()),
	}
	return strings.Join(hashParts, hashPartSeparator), nil
}

// Validate ...
func (proofer ProofOfWork) Validate(block blockchain.Block) error {
	hashParts, targetBit, err := parseHash(block.Hash)
	if err != nil {
		return fmt.Errorf(
			"unable to parse the hash: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	targetBitIndex, err := powValueTypes.NewTargetBitIndex(targetBit)
	if err != nil {
		return fmt.Errorf(
			"unable to construct the target bit index: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	challenge, err := pow.NewChallengeBuilder().
		SetTargetBitIndex(targetBitIndex).
		SetSerializedPayload(powValueTypes.NewSerializedPayload(block.MergedData())).
		SetHash(powValueTypes.NewHash(sha256.New())).
		SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
			"{{ .Challenge.SerializedPayload.ToString }}" +
				"{{ .Nonce.ToString }}" +
				"{{ .Challenge.TargetBitIndex.ToInt }}",
		)).
		Build()
	if err != nil {
		return fmt.Errorf(
			"unable to build the challenge: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	nonce, err := powValueTypes.ParseNonce(hashParts[1])
	if err != nil {
		return fmt.Errorf(
			"unable to parse the nonce: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	rawHashSum, err := hex.DecodeString(hashParts[2])
	if err != nil {
		return fmt.Errorf(
			"unable to decode the hash sum: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	solution, err := pow.NewSolutionBuilder().
		SetChallenge(challenge).
		SetNonce(nonce).
		SetHashSum(powValueTypes.NewHashSum(rawHashSum)).
		Build()
	if err != nil {
		return fmt.Errorf(
			"unable to build the solution: %w",
			errors.Join(err, ErrInvalidParameters),
		)
	}

	if err := solution.Verify(); err != nil {
		return fmt.Errorf("unable to verify the solution: %w", err)
	}

	return nil
}

// Difficulty ...
func (proofer ProofOfWork) Difficulty(hash string) (int, error) {
	_, targetBit, err := parseHash(hash)
	if err != nil {
		return 0, fmt.Errorf("unable to parse the hash: %w", err)
	}

	difficulty := maximalTargetBit - targetBit
	return difficulty, nil
}

func parseHash(hash string) (hashParts []string, targetBit int, err error) {
	hashParts = strings.SplitN(hash, hashPartSeparator, hashPartCount)
	if len(hashParts) != hashPartCount {
		return nil, 0,
			errors.New("the hash contains the invalid quantity of the parts")
	}

	targetBitAsStr := hashParts[0]
	targetBit, err = strconv.Atoi(targetBitAsStr)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to parse the target bit: %w", err)
	}

	return hashParts, targetBit, nil
}
