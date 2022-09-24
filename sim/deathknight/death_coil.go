package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var DeathCoilActionID = core.ActionID{SpellID: 49895}

func (dk *Deathknight) registerDeathCoilSpell() {
	baseDamage := 443.0 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()
	baseCost := float64(core.NewRuneCost(40, 0, 0, 0, 0))
	dk.DeathCoil = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     DeathCoilActionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if dk.SuddenDoomAura.IsActive() {
					cast.GCD = 0
					cast.Cost = 0
				} else {
					cast.GCD = dk.getModifiedGCD()
				}
			},
		},

		BonusCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1.0 + float64(dk.Talents.Morbidity)*0.05) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (baseDamage + 0.15*dk.getImpurityBonus(spell)) * dk.RoRTSBonus(hitEffect.Target)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.OutcomeFuncMagicHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				dk.LastOutcome = spellEffect.Outcome
				if spellEffect.Landed() && dk.Talents.UnholyBlight {
					dk.procUnholyBlight(sim, spellEffect.Target, spellEffect.Damage)
				}
			},
		}),
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 40.0, 0, 0, 0) && dk.DeathCoil.IsReady(sim)
	}, nil)
}

func (dk *Deathknight) registerDrwDeathCoilSpell() {
	baseDamage := 443.0 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()

	dk.RuneWeapon.DeathCoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		BonusCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1.0 + float64(dk.Talents.Morbidity)*0.05) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return baseDamage + 0.15*dk.RuneWeapon.getImpurityBonus(spell)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.RuneWeapon.OutcomeFuncMagicHitAndCrit(),
		}),
	})
}
