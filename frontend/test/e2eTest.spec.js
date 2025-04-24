describe("Augury Frontend Functional Tests", () => {

    it("should load the homepage and display main headings", () => {
      cy.visit("/");
      cy.get("h1").contains("AUGURY").should("be.visible");
      cy.contains("h2", "IOC Intelligence").should("be.visible");
      cy.get('#username').should('be.visible');
      cy.contains("button", "Save Username").should("be.visible");
      cy.contains("p", "Saved Username:").should("be.visible");
      cy.contains("button", "Switch to Light Mode").should("be.visible");
    });
  });
  