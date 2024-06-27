package cli

import (
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	protov1 "github.com/gwuah/rds/api/gen/proto/v1"
	"github.com/gwuah/rds/api/gen/proto/v1/protov1connect"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:          "cli",
		Short:        "Handles cli tasks for rds.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := protov1connect.NewManagerServiceClient(&http.Client{}, "http://0.0.0.0:5555")

			if args[0] == "createDeployment" {
				response, err := client.CreateDeployment(cmd.Context(), &connect.Request[protov1.CreateDeploymentRequest]{
					Msg: &protov1.CreateDeploymentRequest{
						Token:         "FlyV1 fm2_lJPECAAAAAAAAEULxBDgONK5MSPzJZqJd1sNopq4wrVodHRwczovL2FwaS5mbHkuaW8vdjGWAJLOAAJmxB8Lk7lodHRwczovL2FwaS5mbHkuaW8vYWFhL3YxxDwfEViooMd+zquZx4Q7ULj7kBhZFxbqHyQCF4BQCHca1yHCN/acZvj2Ce7CoF/UMHIEoYDQ/FtU2y3J+LvETn1IPGbb0lsn4WqLHjSSCNpB4CP6Oo97JbXM/bKEnuEnw/0qmICFpmxp6rRvjoqrO7PC6U1qjeBOQsG92XCTc6n7GBj68UHAREOr1HkGZQ2SlAORgc4AMP+RHwWRgqdidWlsZGVyH6J3Zx8BxCDAaH83dcKbUHggcP2ZsqIF2MGHULz/eCA0cUgAOlJLQQ==,fm2_lJPETn1IPGbb0lsn4WqLHjSSCNpB4CP6Oo97JbXM/bKEnuEnw/0qmICFpmxp6rRvjoqrO7PC6U1qjeBOQsG92XCTc6n7GBj68UHAREOr1HkGZcQQ0k2hbZkgrqHs+0xkeuZLxsO5aHR0cHM6Ly9hcGkuZmx5LmlvL2FhYS92MZYEks5mfIS0zowUitIKkc4AAjc6DMQQasVfdRniLRIqLbWB7gtn8MQgAUR4YqqdfHfrSVB7OwtdJDB8k2ZZBNA781/rwhSq4PA=",
						AppId:         "cherrypicker-yvonne-90",
						CorrelationId: "4348683643",
						Configs: []*structpb.Struct{
							{
								Fields: map[string]*structpb.Value{},
							},
						},
					},
				})
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(response)
			}
			return nil
		},
	}
}
