package populate

import (
	"bum-service/internal/service/subject"
	"context"
	"fmt"
)

func (s *Service) addSubjects(ctx context.Context) error {
	subjects := []subject.CreateSubjectArgs{
		{
			Name:        "Math",
			Description: pointerFromString("Mathematics, or math, is the study of numbers, quantities, shapes, and spaces. It's a branch of science that uses formulas to relate these concepts."),
		},
		{
			Name:        "Science",
			Description: pointerFromString("Science is the systematic study of the structure and behavior of the physical and natural world through observation and experiment."),
		},
		{
			Name:        "History",
			Description: pointerFromString("History is the study of past events, particularly in human affairs, focusing on civilizations, cultures, and significant milestones."),
		},
		{
			Name:        "Geography",
			Description: pointerFromString("Geography is the study of Earth's landscapes, environments, and the relationships between people and their environments."),
		},
		{
			Name:        "English",
			Description: pointerFromString("English focuses on the study of language, literature, grammar, and communication skills."),
		},
		{
			Name:        "Physics",
			Description: pointerFromString("Physics is the branch of science concerned with the nature and properties of matter and energy, exploring concepts like motion, force, and energy."),
		},
		{
			Name:        "Chemistry",
			Description: pointerFromString("Chemistry is the study of substances, their properties, how they interact, combine, and change to form new substances."),
		},
		{
			Name:        "Biology",
			Description: pointerFromString("Biology is the study of living organisms, their structure, function, growth, evolution, and interactions with their environment."),
		},
		{
			Name:        "Physical Education",
			Description: pointerFromString("Physical Education promotes physical activity and teaches students about health, fitness, and teamwork through sports and exercises."),
		},
		{
			Name:        "Computer Science",
			Description: pointerFromString("Computer Science involves the study of computers, programming, algorithms, and their applications in solving problems."),
		},
		{
			Name:        "Economics",
			Description: pointerFromString("Economics is the study of how societies allocate limited resources to meet their needs and desires, focusing on production, consumption, and trade."),
		},
		{
			Name:        "Civics",
			Description: pointerFromString("Civics is the study of the rights and responsibilities of citizenship, governance, and political systems."),
		},
		{
			Name:        "Foreign Language",
			Description: pointerFromString("Foreign Language involves learning languages other than one's native tongue, focusing on communication skills and cultural understanding."),
		},
		{
			Name:        "Literature",
			Description: pointerFromString("Literature is the study of written works, including prose, poetry, and drama, analyzing themes, styles, and cultural significance."),
		},
		{
			Name:        "Health Education",
			Description: pointerFromString("Health Education focuses on teaching students about personal health, nutrition, hygiene, and well-being."),
		},
		{
			Name:        "Social Studies",
			Description: pointerFromString("Social Studies integrates history, geography, civics, and economics to provide a broad understanding of societies and human behavior."),
		},
		{
			Name:        "Environmental Science",
			Description: pointerFromString("Environmental Science explores the relationships between humans and their environment, focusing on sustainability and conservation."),
		},
	}

	for _, subjectItem := range subjects {
		createdSubject, err := s.subjectService().CreateSubject(ctx, subjectItem)
		if err != nil {
			return fmt.Errorf("failed to create subject: %w", err)
		}

		//s.subjectsList[createdSubject.Name] = createdSubject
		s.subjectsList[createdSubject.ID.String()] = createdSubject
	}

	return nil
}
