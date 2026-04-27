package authz

import (
	"context"
	"io"
	"strings"

	"github.com/MintzyG/FastUtilitiesNet"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

func Subject(kind string, id uuid.UUID) string { return kind + ":" + id.String() }
func Resource(kind string, id string) string   { return kind + ":" + id }
func ResourceType(resType string) string       { return resType }
func Permission(permission string) string      { return permission }

type Caveat struct {
	Name    string
	Context map[string]any
}

func parseRel(s string) (*pb.Relationship, error) {
	hashIdx := strings.Index(s, "#")
	atIdx := strings.Index(s, "@")
	if hashIdx < 0 || atIdx < 0 || atIdx < hashIdx {
		return nil, fun.Errf("invalid rel %q: expected resource:id#relation@subject:id", s).Internal()
	}

	resParts := strings.SplitN(s[:hashIdx], ":", 2)
	subParts := strings.SplitN(s[atIdx+1:], ":", 2)
	relation := s[hashIdx+1 : atIdx]

	if len(resParts) != 2 || len(subParts) != 2 || relation == "" {
		return nil, fun.Errf("invalid rel %q: malformed segments", s).Internal()
	}

	return &pb.Relationship{
		Resource: &pb.ObjectReference{ObjectType: resParts[0], ObjectId: resParts[1]},
		Relation: relation,
		Subject:  &pb.SubjectReference{Object: &pb.ObjectReference{ObjectType: subParts[0], ObjectId: subParts[1]}},
	}, nil
}

func toCaveat(c *Caveat) *pb.ContextualizedCaveat {
	if c == nil {
		return nil
	}
	ctx, _ := structpb.NewStruct(c.Context)
	return &pb.ContextualizedCaveat{CaveatName: c.Name, Context: ctx}
}

func write(ctx context.Context, client *v1.Client, op pb.RelationshipUpdate_Operation, caveat *Caveat, rels ...string) error {
	updates := make([]*pb.RelationshipUpdate, len(rels))
	for i, r := range rels {
		rel, err := parseRel(r)
		if err != nil {
			return err
		}
		rel.OptionalCaveat = toCaveat(caveat)
		updates[i] = &pb.RelationshipUpdate{Operation: op, Relationship: rel}
	}
	_, err := client.WriteRelationships(ctx, &pb.WriteRelationshipsRequest{Updates: updates})
	return err
}

// Can checks "subject:id" has "permission" on "resource:id"
// ex: Can(ctx, client, "user:abc", "edit", "event:xyz")
func Can(ctx context.Context, client *v1.Client, subject, permission, resource string, caveatCtx map[string]any) (bool, error) {
	subType, subID, _ := strings.Cut(subject, ":")
	resType, resID, _ := strings.Cut(resource, ":")

	var pbCtx *structpb.Struct
	if caveatCtx != nil {
		pbCtx, _ = structpb.NewStruct(caveatCtx)
	}

	resp, err := client.CheckPermission(ctx, &pb.CheckPermissionRequest{
		Resource:   &pb.ObjectReference{ObjectType: resType, ObjectId: resID},
		Permission: permission,
		Subject:    &pb.SubjectReference{Object: &pb.ObjectReference{ObjectType: subType, ObjectId: subID}},
		Context:    pbCtx,
	})
	if err != nil {
		return false, fun.Errf("authz check: %v", err).Internal()
	}
	return resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

func Require(ctx context.Context, client *v1.Client, subject, permission, resource string, caveatCtx map[string]any) error {
	allowed, err := Can(ctx, client, subject, permission, resource, caveatCtx)
	if err != nil {
		return err
	}
	if !allowed {
		return fun.ErrForbidden("insufficient permissions")
	}
	return nil
}

// CreateRelation vincula relations, com caveat opcional.
// ex: CreateRelation(ctx, client, "organization:abc#admin@user:xyz")
// ex: CreateRelation(ctx, client, "event:xyz#attendee@user:abc", authz.Caveat{Name: "within_time_range", ...})
func CreateRelation(ctx context.Context, client *v1.Client, rel string, relsAndCaveat ...any) error {
	rels, caveat := splitArgs(rel, relsAndCaveat)
	return write(ctx, client, pb.RelationshipUpdate_OPERATION_TOUCH, caveat, rels...)
}

// DeleteRelation remove relations.
func DeleteRelation(ctx context.Context, client *v1.Client, rel string, rest ...any) error {
	rels, caveat := splitArgs(rel, rest)
	return write(ctx, client, pb.RelationshipUpdate_OPERATION_DELETE, caveat, rels...)
}

// splitArgs separa os strings de rel e o Caveat opcional dos args variádicos.
func splitArgs(first string, rest []any) ([]string, *Caveat) {
	rels := []string{first}
	var caveat *Caveat
	for _, a := range rest {
		switch v := a.(type) {
		case string:
			rels = append(rels, v)
		case Caveat:
			caveat = &v
		}
	}
	return rels, caveat
}

// Lookup retorna os IDs de recursos do tipo resourceType onde o subject tem a permission
// ex: Lookup(ctx, client, "user:abc", "view", "project") -> ["uuid1", "uuid2"]
func Lookup(ctx context.Context, client *v1.Client, subject, permission, resourceType string) ([]string, error) {
	subType, subID, found := strings.Cut(subject, ":")
	if !found {
		return nil, fun.Errf("authz lookup: invalid subject format %q", subject).Internal()
	}
	stream, err := client.LookupResources(ctx, &pb.LookupResourcesRequest{
		ResourceObjectType: resourceType,
		Permission:         permission,
		Subject: &pb.SubjectReference{
			Object: &pb.ObjectReference{
				ObjectType: subType,
				ObjectId:   subID,
			},
		},
	})
	if err != nil {
		return nil, fun.Errf("authz lookup: %v", err).Internal()
	}

	var ids []string
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fun.Errf("authz lookup stream: %v", err).Internal()
		}
		ids = append(ids, resp.ResourceObjectId)
	}
	return ids, nil
}

// Expand retorna todos os subjects que têm a permission no resource
// ex: Expand(ctx, client, "view", "project:abc") -> ["user:uuid1", "user:uuid2"]
func Expand(ctx context.Context, client *v1.Client, permission, resource string) ([]string, error) {
	resType, resID, _ := strings.Cut(resource, ":")

	resp, err := client.ExpandPermissionTree(ctx, &pb.ExpandPermissionTreeRequest{
		Resource:   &pb.ObjectReference{ObjectType: resType, ObjectId: resID},
		Permission: permission,
	})
	if err != nil {
		return nil, fun.Errf("authz expand: %v", err).Internal()
	}

	var subjects []string
	collectLeaves(resp.TreeRoot, &subjects)
	return subjects, nil
}

func collectLeaves(node *pb.PermissionRelationshipTree, out *[]string) {
	if node == nil {
		return
	}
	switch t := node.TreeType.(type) {
	case *pb.PermissionRelationshipTree_Leaf:
		for _, s := range t.Leaf.Subjects {
			*out = append(*out, s.Object.ObjectType+":"+s.Object.ObjectId)
		}
	case *pb.PermissionRelationshipTree_Intermediate:
		for _, child := range t.Intermediate.Children {
			collectLeaves(child, out)
		}
	}
}
