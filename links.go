package cosmos

func createDatabaseLink(databaseID string) string {
	link := "dbs"
	if len(databaseID) > 0 {
		link += "/" + databaseID
	}

	return link
}

func createCollectionLink(databaseID string, collectionID string) string {
	link := createDatabaseLink(databaseID) + "/colls"
	if len(collectionID) > 0 {
		link += "/" + collectionID
	}

	return link
}

func createDocumentLink(databaseID string, collectionID string, documentID string) string {
	link := createCollectionLink(databaseID, collectionID) + "/docs"
	if len(documentID) > 0 {
		link += "/" + documentID
	}

	return link
}

func createAttachmentLink(databaseID string, collectionID string, documentID string, attachmentID string) string {
	link := createDocumentLink(databaseID, collectionID, documentID) + "/attachments"
	if len(attachmentID) > 0 {
		link += "/" + attachmentID
	}

	return link
}
