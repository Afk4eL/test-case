package usecase

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"zametka-service/internal/domain/models"
// 	note_mock "zametka-service/internal/note/mock"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestCreateNote_Success(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	userID := uuid.New().String()
// 	note := &models.Note{
// 		Title: "Test Note",
// 		Data:  "This is a test note",
// 	}
// 	expectedNoteID := uuid.New()

// 	ctx := context.WithValue(context.Background(), "user_id", userID)

// 	mockRepo.On("Create", ctx, mock.Anything).Return(expectedNoteID, nil)

// 	noteID, err := noteUC.CreateNote(ctx, note)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedNoteID, noteID)
// 	assert.Equal(t, note.UserId.String(), userID)
// 	mockRepo.AssertExpectations(t)
// }

// func TestCreateNote_MissingUserID(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	note := &models.Note{
// 		Title: "Test Note",
// 		Data:  "This is a test note",
// 	}

// 	ctx := context.Background()

// 	noteID, err := noteUC.CreateNote(ctx, note)

// 	assert.Error(t, err)
// 	assert.Equal(t, uuid.Nil, noteID)
// 	assert.Equal(t, "No id", err.Error())
// 	mockRepo.AssertNotCalled(t, "Create")
// }

// func TestUpdateNote_Success(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	userID := uuid.New().String()
// 	note := &models.Note{
// 		Id:    uuid.New(),
// 		Title: "Updated Title",
// 		Data:  "Updated Body",
// 	}
// 	expectedUpdatedNote := &models.Note{
// 		Id:    note.Id,
// 		Title: "Updated Title",
// 		Data:  "Updated Body",
// 	}

// 	ctx := context.WithValue(context.Background(), "user_id", userID)

// 	mockRepo.On("Update", ctx, note).Return(expectedUpdatedNote, nil)

// 	updatedNote, err := noteUC.UpdateNote(ctx, note)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUpdatedNote, updatedNote)
// 	mockRepo.AssertExpectations(t)
// }

// func TestDeleteNote_Success(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	noteID := uuid.New()

// 	ctx := context.Background()

// 	mockRepo.On("Delete", ctx, noteID).Return(nil)

// 	err := noteUC.DeleteNote(ctx, noteID)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetAllNotes_Success(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	userID := uuid.New()
// 	expectedNotes := []models.Note{
// 		{Id: uuid.New(), Title: "Note 1", Data: "Body 1"},
// 		{Id: uuid.New(), Title: "Note 2", Data: "Body 2"},
// 	}

// 	ctx := context.Background()

// 	mockRepo.On("GetAll", ctx, userID).Return(expectedNotes, nil)

// 	notes, err := noteUC.GetAllNotes(ctx, userID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedNotes, notes)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetAllNotes_RepoError(t *testing.T) {
// 	mockRepo := new(note_mock.NotePgRepository)
// 	noteUC := NewNoteUsecase(mockRepo)

// 	userID := uuid.New()
// 	expectedError := errors.New("database error")

// 	ctx := context.Background()

// 	mockRepo.On("GetAll", ctx, userID).Return(nil, expectedError)

// 	notes, err := noteUC.GetAllNotes(ctx, userID)

// 	assert.Error(t, err)
// 	assert.Nil(t, notes)
// 	assert.Equal(t, expectedError, err)
// 	mockRepo.AssertExpectations(t)
// }
